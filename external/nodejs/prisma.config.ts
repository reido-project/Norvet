import { defineConfig } from "prisma/config";
import { getDatabaseUrl } from "./env.resolver";

const [databaseUrl, shadowUrl] = getDatabaseUrl();

export default defineConfig({
  schema: "prisma/schema.prisma",
  migrations: {
    path: "prisma/migrations",
    seed: "tsx prisma/seeder/seeder.ts",
  },
  datasource: {
    url: databaseUrl,
    shadowDatabaseUrl: shadowUrl
  },
});
