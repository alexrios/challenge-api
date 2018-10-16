db.swcollection.drop();
db.createCollection("swcollection", {autoIndexId: true});
db.swcollection.createIndex({name: 1}, {unique: true});