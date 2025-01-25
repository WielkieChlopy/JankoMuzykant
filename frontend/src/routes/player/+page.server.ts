import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
    const response = await fetch('http://localhost:8080/api/v1/url', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: 'https://www.youtube.com/watch?v=6v_R180kIGs' }),
    })
    console.log(response)
    const data = await response.json()
    const url = data.url
    return {
        url
    }
};