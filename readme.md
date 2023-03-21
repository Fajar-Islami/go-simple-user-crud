# How to running this app on local
1. copy `.env-example` menjadi `.env`
2. Ubah isi file `.env`   
3. Export env menggunakan `export $(cat .env | xargs -L 1)`  
4. Jika ingin menjalankan di local, isi config docker compose lalu jalankan `docker compose up -d`
5. Lakukan migration dahulu dengan `make migrate-up`
6. Jalankan docker compose aplikasi dengan `make dockerup`
7. Untuk menghentikannya gunakan `make dockerstop`

# TODO
[X] CRUD User   
[X] Get Detail Save and Get From redis   
[X] Connect to 3rd party   
[X] Fix swagger   
[X] Deploy to fly.io   
[X] Add Github Action   
[X] Auto Deploy with GA   
[ ] Migrate from redis to groupcache   
[X] Migrate from viper to godotenv
[X] API HealthCheck   
[X] Fix swagger port, change host to different with normal host