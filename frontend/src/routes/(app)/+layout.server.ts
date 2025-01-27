import type { LayoutServerLoad } from './$types'
import { fail, redirect } from '@sveltejs/kit'

export const load: LayoutServerLoad = async (event) => {
	const user = await event.locals.user
	if (!user) return redirect(303, `/auth?redirect=${event.url.pathname}`)
}
