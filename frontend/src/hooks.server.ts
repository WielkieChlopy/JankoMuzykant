import { sequence } from '@sveltejs/kit/hooks'
import type { Handle } from '@sveltejs/kit';
import { i18n } from '$lib/i18n';
import jwt from 'jsonwebtoken';
import { JWT_SECRET } from '$env/static/private';

const handleParaglide: Handle = i18n.handle();

const handleAuth: Handle = async ({ event, resolve }) => {
	const token = event.cookies.get('token');
	if (!token) {
		return resolve(event);
	}

	let decoded;
	try {
		decoded = jwt.verify(token, JWT_SECRET);

		if (typeof decoded != 'object' || decoded === null) {;
			return resolve(event);
		}
	} catch (error) {
		return resolve(event);
	}

	const { id, username, exp } = decoded as { id: string; username: string; exp: number };
	event.locals.user = { id, username, exp };
	event.locals.token = token;

	return resolve(event);
}

export const handle: Handle = sequence(handleParaglide, handleAuth);
