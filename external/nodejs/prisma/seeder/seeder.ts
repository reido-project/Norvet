import { PrismaClient } from "./_prisma/client";
import { getPrismaAdapter } from "../../env.resolver";
import * as seed from "./index";

export const prisma = new PrismaClient({ adapter: getPrismaAdapter() as any });

async function main() {
    console.log("🌱 Starting database seed...\n");

    await seed.users(prisma);
}

main()
    .then(async () => {
        console.log("🌱 Database seeded successfully.");
        await prisma.$disconnect();
    })
    .catch(async (err) => {
        console.error("❌ Error seeding database:", err);
        await prisma.$disconnect();
        process.exit(1);
    });