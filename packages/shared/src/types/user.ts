/**
 * Profil aplikasi yang melengkapi `auth.users` Supabase. `id` sengaja sama persis
 * dengan `auth.users.id` (relasi satu-ke-satu), bukan primary key sendiri, supaya
 * tidak ada dua sumber kebenaran soal "siapa user ini".
 */
export interface User {
	id: string;
	name: string;
	phone: string | null;
	/**
	 * Asal provinsi, disimpan sebagai kode stabil dari INDONESIA_PROVINCES
	 * (mis. `jawa-barat`), bukan nama tampilannya. Nullable karena akun yang
	 * dibuat lewat Google OAuth belum pernah melewati form yang menanyakannya.
	 */
	province: string | null;
	avatarUrl: string | null;
	createdAt: string;
}
