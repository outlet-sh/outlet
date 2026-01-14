import pg from 'pg';
import { env } from '$env/dynamic/private';
const { DATABASE_URL } = env;

const { Pool } = pg;

// Create a connection pool
export const pool = new Pool({
	connectionString: DATABASE_URL,
	max: 20,
	idleTimeoutMillis: 30000,
	connectionTimeoutMillis: 2000
});

// Helper function to execute queries
export async function query(text: string, params?: any[]) {
	const start = Date.now();
	try {
		const res = await pool.query(text, params);
		const duration = Date.now() - start;
		console.log('Executed query', { text, duration, rows: res.rowCount });
		return res;
	} catch (error) {
		console.error('Database query error:', error);
		throw error;
	}
}

// Graceful shutdown
process.on('SIGINT', async () => {
	await pool.end();
	process.exit(0);
});

process.on('SIGTERM', async () => {
	await pool.end();
	process.exit(0);
});
