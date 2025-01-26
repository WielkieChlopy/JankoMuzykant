import { BACKEND_URL } from '$env/static/private';
import { fail, redirect } from '@sveltejs/kit';

export const actions = {
	login: async ({ cookies, request, url }) => {
		const formData = await request.formData();
		const email = formData.get('email');
		const password = formData.get('password');
		const redirectPath = url.searchParams.get('redirect') || '/home';

		const response = await fetch(`${BACKEND_URL}/auth/login`, {
			method: 'POST',
			body: JSON.stringify({ email, password }),
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

    logout: async ({ cookies }) => {
        cookies.delete('token', { path: '/' });
        return { success: true };
    }
};
