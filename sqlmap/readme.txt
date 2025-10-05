sqlmap -r req.txt  --batch -v 3
sqlmap -r ds.txt  --batch -v 3 --risk=3 --level=5
sqlmap -r ds.txt  --batch  --random-agent  --level=5 --risk=3  --threads=10  --timeout=10 --retries=2   --technique=BEUSTQ  --tamper=space2comment,randomcase,charencode -p "filter" 
sqlmap -r ds.txt --batch --random-agent --technique=U --union-cols=1-5 -p filter --threads=3 --timeout=15
