
publisher.py - secure connection
---------------------------------
python publisher.py -v --pull tcp://127.0.0.1:7000 --pub tcp://127.0.0.1:7001



publisher.py - insecure connection
----------------------------------
python publisher.py -v --insecure --pull tcp://127.0.0.1:7000 --pub tcp://127.0.0.1:7001



cache.py - (No option for security)
-----------------------------------------
python cache.py -v --sub-address tcp://127.0.0.1:7001 --router-address tcp://127.0.01:8000


soundtouch
----------
./soundtouch 


collector
----------
./collector