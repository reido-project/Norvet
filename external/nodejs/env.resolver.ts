import dotenv from "dotenv";
import path from "path";
import { PrismaPg } from "@prisma/adapter-pg";
import { PrismaMariaDb } from "@prisma/adapter-mariadb";

dotenv.config({ path: path.resolve(__dirname, "../../.env") });

export function getDatabaseUrl(): [string, string] {
    const driver = process.env["DB_DRIVER"]?.toLowerCase();
    const host = process.env["DB_HOST"] ?? "localhost";
    const port = process.env["DB_PORT"] ?? "3306";
    const user = process.env["DB_USER"] ?? "";
    const password = process.env["DB_PASSWORD"] ?? "";
    const name = process.env["DB_NAME"] ?? "";
    const shadowName = process.env["SHADOW_DATABASE"] ?? "nub1_shadow";

    let dbUrl: string;
    let shadowUrl: string;

    if (driver === "postgres" || driver === "postgresql") {
        const sslmode = process.env["POSTGRES_SSLMODE"] ?? "disable";
        const timezone = process.env["POSTGRES_TIMEZONE"] ?? "UTC";

        dbUrl = `postgresql://${user}:${password}@${host}:${port}/${name}?sslmode=${sslmode}&timezone=${timezone}`;
        shadowUrl = `postgresql://${user}:${password}@${host}:${port}/${shadowName}?sslmode=${sslmode}&timezone=${timezone}`;
    } else {
        const charset = process.env["MARIADB_CHARSET"] ?? "utf8mb4";
        const loc = process.env["MARIADB_LOC"] ?? "UTC";
        const encodedLoc = encodeURIComponent(loc);

        dbUrl = `mysql://${user}:${password}@${host}:${port}/${name}?charset=${charset}&loc=${encodedLoc}`;
        shadowUrl = `mysql://${user}:${password}@${host}:${port}/${shadowName}?charset=${charset}&loc=${encodedLoc}`;
    }

    return [dbUrl, shadowUrl];
}

export function getPrismaAdapter() {
    const driver = process.env["DB_DRIVER"]?.toLowerCase();
    const [databaseUrl] = getDatabaseUrl(); // Memanfaatkan fungsi di atas

    console.log();

    if (driver === "postgres" || driver === "postgresql") {
        console.log("Using database: PostgreSQL");
        return new PrismaPg({ connectionString: databaseUrl });
    }

    console.log("Using database: MySQL / MariaDB");
    console.log();
    return new PrismaMariaDb(databaseUrl);
}