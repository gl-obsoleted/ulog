
set DEV_DB=db_dev_test

if not exist %DEV_DB% md %DEV_DB%

tools\mongodb-2.6.5\mongod.exe --dbpath "%~dp0%DEV_DB%"
