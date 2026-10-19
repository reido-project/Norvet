import { PrismaClient } from "../_prisma/client";
import { hashPassword } from "../libs/bcrypt";

export async function users(prisma: PrismaClient){
    console.log("↻ Seeding users...");

    const users = [
        { username: "ruddiularsawah", email: "ruddiulsaw@mail.com", password: "ambatoukum" }
    ];

    for(const user of users){
        const hashedPassword = await hashPassword(user.password)

        await prisma.user.upsert({
            where: { email: user.email },
            update: {
                username: user.username,
                password: hashedPassword
            },
            create: {
                username: user.username,
                email: user.email,
                password: hashedPassword
            }
        });
    }

    console.log("✔ Table 'users' seeded successfully.\n");
}