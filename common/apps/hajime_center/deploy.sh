sudo apt-get update
sudo apt-get install supervisor

sudo mkdir /srv/HajimeCenter
sudo mkdir /srv/HajimeCenter/logs
go build -o /srv/HajimeCenter/HajimeCenter main.go

cp HajimeCenter.conf /etc/supervisor/conf.d/
sudo supervisorctl reread
sudo supervisorctl update
sudo supervisorctl start HajimeCenter
