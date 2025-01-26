import { add_song, get_queue } from '@/components/server/queue';
import type { PageServerLoad } from './$types';
import { BACKEND_URL } from '$env/static/private';

export const load: PageServerLoad = async () => {
    return {
        queue: get_queue()
    }
};

export const actions = {
    add_song: async ({ request }) => {
        const formData = await request.formData();
        const url = formData.get('url');

        const response = await fetch(`${BACKEND_URL}/api/v1/details`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ url: url }),
        })
        const data = await response.json()
        add_song(data)
    }
}