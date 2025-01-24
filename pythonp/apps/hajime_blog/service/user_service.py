import logging

from passlib.context import CryptContext

from pythonp.apps.hajime_blog.db.models import BizUser, User, UserStat
from pythonp.apps.hajime_blog.db.schemas import BizUserModel, BizLoginModel, BizAuthResponseModel, SignLoginModel, \
    UserAuthResponseModel
from pythonp.apps.hajime_blog.utils.common import error_return, success_return, get_session_id
from pythonp.apps.hajime_blog.utils.redis import get_redis

class UserService:
    def __init__(self):
        super().__init__()
        self.hasher = CryptContext(schemes=['bcrypt'], deprecated='auto')
        self.secret_key = '0xbotxdataloveyou'
        self.algorithm = 'HS256'
        self.access_token_expire_minutes = 60 * 24 * 30
        self.guard = {}


    def clear_guard(self):
        self.guard = {}
    async def create_biz_user(self, form: BizUserModel):
        username = form.username
        password = form.password

        print(form)
        user = await BizUser.find_one({"username": username})
        if user is not None:
            return error_return(1, "username already exists", user)
        password = self.hasher.hash(password)
        user = await  BizUser(
            username=username,
            password=password,
            ga=form.ga
        ).create()
        return success_return(user)


    async def biz_logout(self,oid):
        redis = await get_redis()
        await redis.delete(oid)
    async def biz_user_login(self, login_form: BizLoginModel):
        username = login_form.username
        password = login_form.password

        biz_user = await BizUser.find_one(BizUser.username == username)

        if biz_user is None:
            return error_return(418, "username doesn't exists")

        if self.hasher.verify(password, biz_user.password):
            access_token = await self.create_access_token(biz_user.id)
            data = {
                "token": access_token,
                "uid": biz_user.id,
                "acl": ["all"],
                "username": username,
                "avatar": biz_user.get_avatar()
            }
            auth = BizAuthResponseModel(**data)
            return success_return(auth)

        else:
            return error_return(418, "Password Error")

    async def create_access_token(self, uid: int,loc:str="front")->str:
        session_id = get_session_id()

        logging.info("create_access_token :%s ", session_id)
        uid_key = loc+":"+str(uid)
        redis = await get_redis()

        old_session_id = await redis.get(uid_key)
        if old_session_id:
            logging.info("delete old session_id :%s ", old_session_id)
            # await redis.delete(old_session_id)
        await redis.set(session_id, uid)
        await redis.set(uid_key, session_id)
        await redis.expire(session_id, 60 * self.access_token_expire_minutes)

        return session_id

    async def logout(self, uid):
        logging.info("logout:{}".format(uid))
        redis = await get_redis()
        session_id = await redis.get(uid)
        if session_id is not None:
            await redis.delete(session_id)
        sid = await redis.get(f"socket:{uid}")
        if sid is not None:
            await redis.delete(f"socket:{uid}")
            await redis.delete(f"socket:{sid}")
            await redis.delete(f"socket:{sid}:token")
            logging.info("logout socket:{}".format(sid))
        return success_return()


    async def add_refer(self,uid,code):
        print("add_refer",uid,code)
        ref_user = await User.get_user_by_code(code)
        if ref_user is not None and ref_user.id == uid:
            print("can not add self")
            return error_return(1, "can not add self")
        if ref_user is not None:
            pid = ref_user.id
            ids = await UserStat.get_parents(pid, 10000)
            if uid not in ids:
                # set parent
                print("set parent")
                await User.find_one(User.id == uid).update_one({"$set": {"pid": pid,"p_address":ref_user.address}})
                print("set parent")
                await UserStat.find_one({"uid": uid}).update_one({"$set": {"pid": pid,"p_address":ref_user.address}})
                # set counter
                await UserStat.find_one({"uid": pid}).update_one({"$inc": {"direct_child": 1}})
                # set counter
                parent_ids = await UserStat.get_parents(uid, 10000)
                userstat = await UserStat.find_one({"uid": uid})
                total_buy_amount = userstat.total_buy_amount

                await UserStat.find_one({"uid": {"$in": parent_ids}}).update_one(
                    {"$inc": {"total_team_num": 1, "total_buy_amount": total_buy_amount}})
                print("add_refer_succ")
                return success_return()
            else:
                print("循环邀请错误")

                return error_return(1, "循环邀请错误")
        else:
            print("邀请码错误")
            return error_return(1, "邀请码错误")

    async def login_with_sign(self, form:SignLoginModel):
        address = form.walletAddress.lower()

        user = await User.get_user_by_address(address)
        if user is None:
            user = await User(address=address).create()
            if len(form.code)>2:
                await self.add_refer(user.id, form.code)


        if user:
            access_token = await self.create_access_token(user.id)
            # parent_address = ""
            # if user.pid:
            #     parent_user = await User.get_user_by_uid(user.pid)
            #     if parent_user is not None:
            #         parent_address = parent_user.address
            data = {
                "token": access_token,
                "uid": user.id,
                "code": user.code,
                "invite_link": user.get_invite_link(),
                "pid": user.pid,
                "parent":user.p_address
            }
            redis = await get_redis()
            await redis.setex(access_token, 60 * self.access_token_expire_minutes, user.id )
            auth =  UserAuthResponseModel(**data)
            return success_return(auth)
        else:
            return error_return(418, "Wallet Address Error")

    async def get_all_childs(self, uid):

        cursor = await UserStat.find({'pid': uid}).to_list()

        children = []
        i = 0
        for o in cursor:
            stat =  await UserStat.find_one({'uid': o.uid})
            buy_amount = stat.buy_amount
            total_buy_amount = stat.total_buy_amount
            t = await User.get_user_by_uid(o.uid)

            item = {
                'label': f"[{t.address},Individual:{buy_amount}-Team:{total_buy_amount}]",
            }
            i += 1
            if o.uid in self.guard:
                pass

            vo = await self.get_all_childs(o.uid)
            item['children'] = vo[0]
            item['child_cnt'] = vo[1]
            i += item['child_cnt']
            self.guard[o.uid] = o.uid
            children.append(item)
        return [children,i]
    async def get_childs(self,uid):
        self.guard = {}
        data = await self.get_all_childs(uid)
        out = {
            'nodes':data[0],
            'num':data[1]
        }
        return out
