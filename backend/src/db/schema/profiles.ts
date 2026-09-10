import { pgTable, text, timestamp, uuid } from 'drizzle-orm/pg-core';
import { authUsers } from './_supabase-auth';

/**
 * Profil aplikasi yang melengkapi `auth.users`. `id` sengaja sama dengan
 * `auth.users.id` (relasi satu-ke-satu), bukan primary key sendiri.
 */
export const profiles = pgTable('profiles', {
	id: uuid('id')
		.primaryKey()
		.references(() => authUsers.id, { onDelete: 'cascade' }),
	name: text('name').notNull(),
	phone: text('phone'),
	/**
	 * Asal provinsi, diisi saat pendaftaran. Menyimpan KODE dari
	 * INDONESIA_PROVINCES (mis. `jawa-barat`), bukan nama tampilannya — lihat
	 * alasannya di packages/shared/src/constants/provinces.ts.
	 *
	 * Sengaja `text` biasa, bukan pgEnum seperti event_type/vendor_category:
	 * daftar provinsi ditentukan pemerintah dan memang berubah (pemekaran Papua
	 * 2022 menambah 4 provinsi sekaligus). Kolom enum berarti setiap pemekaran
	 * butuh migrasi ALTER TYPE, padahal tidak ada logika database yang bergantung
	 * pada nilai kolom ini. Validasi nilainya cukup di provinceSchema.
	 *
	 * Nullable: akun lama dan akun yang mendaftar lewat Google OAuth belum pernah
	 * melewati form yang menanyakan asal.
	 */
	province: text('province'),
	avatarUrl: text('avatar_url'),
	createdAt: timestamp('created_at', { withTimezone: true }).notNull().defaultNow(),
});
