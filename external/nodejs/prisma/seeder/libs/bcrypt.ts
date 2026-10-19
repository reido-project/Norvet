import bcrypt from "bcryptjs";

const SALT_ROUNDS = 10;

export async function hashPassword(plainPassword: string): Promise<string> {
    if (!plainPassword) return "";
    return await bcrypt.hash(plainPassword, SALT_ROUNDS);
}