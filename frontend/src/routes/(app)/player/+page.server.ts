import type { PageServerLoad } from './$types';
import { BACKEND_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ locals }) => {
    return {
        backend_url: BACKEND_URL,
        token: locals.token
    }
};

export const actions = {
    add_song: async ({ request }) => {
        const formData = await request.formData();
        const url = formData.get('url');

        const response = await fetch(`${BACKEND_URL}/api/v1/songs/details`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: url }),
        })
        const data = await response.json()
    }
}