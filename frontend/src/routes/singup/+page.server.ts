import { BACKEND_URL } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';

export const actions = {
	signup: async ({ cookies, request, url }) => {
		const formData = await request.formData();
		const username = formData.get('username');
		const password = formData.get('password');
		const redirectPath = url.searchParams.get('redirect') || '/player';

		const response = await fetch(`${BACKEND_URL}/api/v1/signup`, {
			method: 'POST',
			body: JSON.stringify({ username, password }),
			headers: {
				'Content-Type': 'application/json',
			},
		});

		if (response.ok) {
			const data = await response.json();
			cookies.set('token', data.token, {
				path: '/',
				// httpOnly: true,
				// secure: true,
				// sameSite: 'strict', //todo:
			});
			return redirect(303, redirectPath);
		} else {
			return fail(401, { error: 'Invalid credentials' });
		}
	},
};
