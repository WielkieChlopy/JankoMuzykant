import type { PageServerLoad } from './$types';
import { BACKEND_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ locals }) => {
    const response = await fetch(`${BACKEND_URL}/api/v1/queue/`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${locals.token}`
        },
    })
    const data = await response.json()
    console.log("queue", data)
    return {
        backend_url: BACKEND_URL,
        token: locals.token,
        queue: data
    }
};

export const actions = {
    add_song: async ({ request, locals }) => {
        const formData = await request.formData();
        const url = formData.get('url');

        const response = await fetch(`${BACKEND_URL}/api/v1/queue/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${locals.token}`
            },
            body: JSON.stringify({ url: url }),
        })
        const data = await response.json()
    }
}