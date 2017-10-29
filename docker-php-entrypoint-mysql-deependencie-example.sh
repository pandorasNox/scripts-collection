#!/bin/sh

echo ""

echo "START run entrypoint"
echo ""

# bash ./check_mysql_connection.sh
# exit_status=$?
# if [ $exit_status -eq 1 ]; then
#     echo "    exit_status was 1 (error)"
# fi
# if [ $exit_status -eq 0 ]; then
#     echo "    exit_status was 0 (success)"
# fi

    # DBCONNECTION=$(mysql -h $MYSQL_HOST -P $MYSQL_PORT -u$MYSQL_USER -p$MYSQL_PASSWORD --batch --skip-column-names -e "SHOW DATABASES;" 2> /dev/null )
    # if [ -n "$DBCONNECTION" ]
    # then
    #     echo "Connection ok"
    # fi

    # if [ ! -n "$DBCONNECTION" ]
    # then
    #     echo "Connection error"
    # fi

    # exit

#bash check_mysql_connection.sh

resultCheckMysqlConnection=$(bash -c '(./check_mysql_connection.sh); exit $?' 2>&1)
exit_status=$?
if [ $exit_status -eq 1 ]; then
    echo "    exit_status was 1 (error)"
    echo "    reason: $resultCheckMysqlConnection"
    echo ""
fi
if [ $exit_status -eq 0 ]; then
    echo "    exit_status was 0 (success)"
    echo ""
fi

echo start php-fpm
php-fpm

echo "END run entrypoint"
echo ""

# Capturing output and exit codes in BASH / SHELL
# resultA=$(bash -c '(./a.out); exit $?' 2>&1)
# exitA=$?
# resultB=$(bash -c '(./b.out); exit $?' 2>&1)
# exitB=$?


#todo: check db connection with 5 retries




# php /temp/env.php

# output=`php -f somescript.php`
# exitcode=$?

# anotheroutput=`php -f anotherscript.php`
# anotherexitcode=$?

# EMAIL="myemail@foo.com"
# PATH=/sbin:/usr/sbin:/usr/bin:/usr/local/bin:/bin
# export PATH

# output=`php-cgi -f /www/Web/myscript.php myUrlParam=1`
# #echo $output

# if [ "$output" = "0" ]; then
#    echo "Success :D"
# fi
# if [ "$output" = "1" ]; then
#    echo "Failure D:"
#    mailx -s "Script failed" $EMAIL <<!EOF
#      This is an automated message. The script failed.

#      Output was:
#        $output
# !EOF
# fi




#EXPOSE 9000
#php-fpm
