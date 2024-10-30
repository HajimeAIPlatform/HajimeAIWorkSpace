from typing import TypeVar, Generic, List, Any

from pydantic import BaseModel, Field, EmailStr

T = TypeVar('T')


class GenericResponseModel(BaseModel, Generic[T]):
    code: int = 0
    message: str = ""
    result: T = None


class EvidenceDataModel(BaseModel):
    task_id: str
    data: str
    data_hash: str


class EvidenceDataInputModel(BaseModel):
    task_id: str
    data: str
    node_id: str = ""


class TaskQueryModel(BaseModel):
    task_id: str


class UserDepositModel(BaseModel):
    id: str = Field(..., alias='_id')
    address: str
    signature: List[str]
    changeType: str
    changeAmount: int
    decimals: int
    postBalance: int
    preBalance: int
    tokenAddress: str
    owner: str
    blockTime: int
    slot: int
    fee: int
    symbol: str


class SolScanModel(BaseModel):
    total: int
    data: List[UserDepositModel]


class SignLoginModel(BaseModel):
    walletAddress: str = Field(min_length=30, max_length=80)
    sign: str
    code: str = ""
    msg: str = ""


class BizLoginModel(BaseModel):
    username: str
    password: str
    code: str = Field("", description="GA Code")
    sid: str = ""


class BizUserModel(BaseModel):
    username: str = ""
    password: str = ""
    ga: str = ""
    avatar: str = ""


class UserAuthResponseModel(BaseModel):
    uid: str
    pid: str = ""
    parent: str = ""
    token: str
    code: str
    invite_link: str = ""


class BizAuthResponseModel(BaseModel):
    uid: str
    token: str
    username: str = ""
    avatar: str = ""
    acl: List[str] = []


class UserStateModel(BaseModel):
    uid: str = ""
    pid: str = ""


class ReferModel(BaseModel):
    code: str


class NewsSubscribeModel(BaseModel):
    email: EmailStr


class IDModel(BaseModel):
    id: str


class UserInfoModel(BaseModel):
    id: str
    address: str = ""
    p_address: str = ""
    code: str = ""


class UserListRequestModel(BaseModel):
    page: int = 1
    pagesize: int = 10
    address: str = ""


class UserListItemModel(BaseModel):
    id: str
    address: str
    p_address: str
    code: str
    create_at: int


class UserListResponseModel(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[UserListItemModel] = []


class UserRewardListItemModel(BaseModel):
    uid: str = ""
    amount: int = 0
    description: int = 0
    create_at: int = 0


class UserRewardListResponse(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[UserRewardListItemModel] = []


class GeneralPageRequest(BaseModel):
    page: int = 1
    pagesize:int=20



class UserInviteListItemRequest(BaseModel):
    page: int = 1


class UserInviteListItemModel(BaseModel):
    uid: str
    address: str
    buy_amount: int = 0
    total_buy_amount: int = 0


class UserInviteListResponse(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[UserInviteListItemModel] = []


class UserStatResponse(BaseModel):
    buy_amount: int = 0
    total_buy_amount: int = 0
    total_team_num: int = 0
    direct_buy_amount: int = 0


class MinerOrderRequest(BaseModel):
    address: str = Field(min_length=0, max_length=80)
    tx_hash: str = Field(min_length=0, max_length=180)
    pay_method: str
    product_id: int
    price:str = ""


class MinerListItemModel(BaseModel):
    id: str
    address: str
    tx_hash: str
    nft_type: str=""
    pay_status: int=1
    amount:float=0
    create_at: int


class MinerListResponse(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[MinerListItemModel] = []


### for admin ###
class AdminUserListItemModel(BaseModel):
    id: str
    uid: str
    address: str = ""
    p_address: str = ""
    code: str = ""
    buy_amount: int = 0
    total_buy_amount: int = 0
    direct_child: int = 0
    total_team_num: int = 0

    create_at: int


class AdminUserListResponseModel(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[AdminUserListItemModel] = []


class AdminDepositListItemModel(BaseModel):
    id: str
    uid: str = ""
    address: str = ""
    signature: List[str]
    changeType: str
    changeAmount: int
    decimals: int
    postBalance: int
    preBalance: int
    tokenAddress: str
    owner: str
    blockTime: int
    slot: int
    fee: int
    symbol: str
    sender: str = ""
    nft_type: str = ""
    create_at: int


class AdminDepositListResponseModel(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    total_deposit:float=0
    list: List[AdminDepositListItemModel] = []


class AdminDepositListRequestModel(BaseModel):
    username: str = ""
    start_date:str=""
    end_date:str=""
    type: str = ""
    page:int=0
    pagesize:int=100

class UIDModel(BaseModel):
    uid: str


##

class APIDepositListItemModel(BaseModel):
    id: str
    address: str = ""
    signature: List[str]
    changeType: str
    changeAmount: int
    tokenAddress: str
    owner: str
    blockTime: int
    symbol: str
    sender: str = ""
    nft_type: str = ""
    create_at: int


class APIDepositListResponseModel(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 10
    list: List[APIDepositListItemModel] = []


class AdminOplogListItemModel(BaseModel):
    id: str
    msg:  str = ""
    uid:  str = ""
    address:str=""
    extra:  str = ""
    op: str = ""
    create_at: int
class AdminOplogListResponseModel(BaseModel):
    total: int = 0
    total_page: int = 1
    pagesize: int = 20
    list: List[AdminOplogListItemModel] = []

class AdminOplogListRequestModel(BaseModel):
    username: str = ""
    page:int=0
    pagesize:int=100

class AdminBuyGroup(BaseModel):
    id: str = Field(None, alias="_id")
    total: float



class UploadResponseModel(BaseModel):
    id:  str = ""
    url: str = ""
    size: int = 0