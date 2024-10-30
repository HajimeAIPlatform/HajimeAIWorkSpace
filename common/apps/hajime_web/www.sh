#!/bin/bash
git commit -m "激活变更" -a
rsync -avz --exclude-from='exclude-file.txt' -e 'ssh -i /Users/mac/.ssh/hajime-website.pem' /Users/mac/workspace/solona/hajime-web ubuntu@18.179.12.70:/opt/hajime

