build:
	go build -o candy-client main.go 
	go build -C ../server -o candy-server main.go

clean:
	@rm -rf log1.txt log2.txt log3.txt log4.txt log5.txt true_log1.txt true_log2.txt true_log3.txt true_log4.txt true_log5.txt ../server/candy-server candy-client

test: test1 test2 test3 test4 test5

test1:
	@touch true_log1.txt
	@echo 'Thank you! Your change is 0' > true_log1.txt
	./candy-client -k YR -c 2 -m 46 > log1.txt
	diff true_log1.txt log1.txt

test2:
	@touch true_log2.txt
	@echo 'Thank you! Your change is 5' > true_log2.txt
	./candy-client -k AA -c 1 -m 20 > log2.txt
	diff true_log2.txt log2.txt

test3:
	@touch true_log3.txt
	@echo 'error: Invalid candy type' > true_log3.txt
	./candy-client -k sweet -c 2 -m 46 > log3.txt
	diff true_log3.txt log3.txt

test4:
	@touch true_log4.txt
	@echo 'error: some error in input numerical values' > true_log4.txt
	./candy-client -k AA -c 2 -m -46 > log4.txt
	diff true_log4.txt log4.txt

test5:
	@touch true_log5.txt
	@echo 'error: You need 30 money!' > true_log5.txt
	./candy-client -k AA -c 2 -m 15 > log5.txt
	diff true_log5.txt log5.txt