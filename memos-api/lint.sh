find . | grep  '.*.py$' |grep -v ".env" |grep -v "history" |xargs -I {} autoflake --in-place --remove-unused-variables {}
find . | grep  '.*.py$' |grep -v ".env" |grep -v "history" |xargs yapf -i #格式化
pipreqs . --encoding=utf-8 --force --ignore .env,.history,__pycache__,venv