use orders

db.orders.createIndex({orderId: 1}, {unique: true})