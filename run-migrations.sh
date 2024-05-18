# run-migrations.sh
#!/bin/sh
for f in /docker-entrypoint-initdb.d/*.sql; do
    echo "Running $f"
    mysql -h mysql -uroot -proot orders < "$f"
done