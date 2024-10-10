console.log('Starting MongoDB setup...');

const env = process.env;

db.createUser({
  user: env.DB_USER,
  pwd: env.DB_PASS,
  roles: [{ role: 'readWrite', db: env.DB_NAME }],
});

db.createCollection('users');

console.log('MongoDB setup finished successfully!');
