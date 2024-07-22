-- DropForeignKey
ALTER TABLE "invites" DROP CONSTRAINT "invites_author_id_fkey";

-- AddForeignKey
ALTER TABLE "invites" ADD CONSTRAINT "invites_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "users"("id") ON DELETE SET NULL ON UPDATE CASCADE;
