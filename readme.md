# How to running this app on local
1. copy `.env-example` menjadi `.env`
2. Isi file `.env`
3. Jika ingin menjalankan di local gunakan `docker compose up -d`
4. Lakukan migration dahulu dengan `migrate-up`
5. Build aplikasi menggunakan `make build`
6. Untuk menjalan gunakan `make run`

# TODO
[X] CRUD User   
[X] Get Detail Save and Get From redis   
[X] Connect to 3rd party   
[X] Fix swagger   
[ ] Deploy to fly.io   
[X] Add Github Action   
[ ] Auto Deploy with GA   
[ ] Migrate from redis to groupcache
[ ] Migrate from viper to github.com/mashingan/smapping
[ ] API HealthCheck