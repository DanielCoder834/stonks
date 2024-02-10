import { createClient } from 'redis';
import dotenv from 'dotenv';
dotenv.config();
export const client = createClient({
    // password: process.env.REDIS_PASSWORD,
    //process.env.REDIS_HOST,
    //parseInt(process.env.REDIS_PORT || '6379', 10)
    socket: {
        host: "127.0.0.1",
        port: 6379
    }
});
