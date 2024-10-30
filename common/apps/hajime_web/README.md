# Hajime-evidence
The Hajime-evidence is designed to enhance the efficiency of proof of work and the security of data. 
By receiving task-related data from Hajimebot nodes, such as vector embedding, vector querying, speech synthesis, and large model inference, it can effectively process and store a vast amount of data, and securely record the hash values of this data on the Solana blockchain, ensuring the data's immutability and transparency.

Please note that not all features have been implemented during the POC phase.






## Prerequisites
- python 3.10
- fastapi
- mongodb


# install
```
poetry env use 3.10
poetry shell
poetry install

python hajime.py
```




curl 'https://www.hajime.ai/admin/biz_login' \
  -H 'Accept: application/json, text/plain, */*' \
  -H 'Accept-Language: en' \
  -H 'Cache-Control: no-cache' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'Origin: https://www.hajime.ai' \
  -H 'Pragma: no-cache' \
  -H 'Referer: https://www.hajime.ai/adminblog/' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  -H 'User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1' \
  -H 'X-Requested-With: XMLHttpRequest' \
  --data-raw '{"username":"admin","password":"admin888"}'


curl -X 'POST' \
  'https://www.hajime.ai/v2/api/market/can_buy_miner' \
  -H 'accept: application/json' \
  -d ''