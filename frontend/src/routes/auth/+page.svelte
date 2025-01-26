<script lang="ts">
	import { Button } from "$lib/components/ui/button/index.js";
	import { Icons } from "$lib/components/icons/icons.js";
	import { Input } from "$lib/components/ui/input/index.js";
	import { Label } from "$lib/components/ui/label/index.js";
	import { cn } from "$lib/utils.js";
	import { enhance } from '$app/forms';
	import * as m from '$lib/paraglide/messages.js';

	let className: string | undefined | null = undefined;
	export { className as class };

	let isLoading = false;
</script>

<div
	class="container relative hidden h-[800px] flex-col items-center justify-center md:grid lg:max-w-none  lg:px-0"
>
	<div class="lg:p-8">
		<div class="mx-auto flex w-full flex-col justify-center space-y-6 sm:w-[350px]">
			<div class="flex flex-col space-y-2 text-center">
				<h1 class="text-2xl font-semibold tracking-tight">{m.login()}</h1>
				<p class="text-muted-foreground text-sm">
					{m.enter_email()}
				</p>
			</div>
			
			<div class={cn("grid gap-6", className)} {...$$restProps}>
				<form 
					method="POST"
					action="?/login" 
					use:enhance={() => {
					isLoading = true;
			
						return async ({ update }) => {
							await update();
							isLoading = false;
						};
					}}
				>
					<div class="grid gap-2">
						<div class="grid gap-1">
							<Label class="sr-only" for="email">{m.email()}</Label>
							<Input
								id="email"
								name="email"
								placeholder="name@example.com"
								type="email"
								autocapitalize="none"
								autocomplete="email"
								autocorrect="off"
								disabled={isLoading}
							/>
							<Label class="sr-only" for="password">{m.password()}</Label>
							<Input
								id="password"
								name="password"
								placeholder="********"
								type="password"
								disabled={isLoading}
							/>
						</div>
						<Button type="submit" disabled={isLoading}>
							{#if isLoading}
								<Icons.spinner class="mr-2 h-4 w-4 animate-spin" />
							{/if}
							{m.login()}
						</Button>
					</div>
				</form>
				<div class="relative">
					<div class="absolute inset-0 flex items-center">
						<span class="w-full border-t" > </span>
					</div>
					<div class="relative flex justify-center text-xs uppercase">
						<span class="bg-background text-muted-foreground px-2">{m.or_continue_with()}</span>
					</div>
				</div>
				<Button variant="outline" type="button" disabled={isLoading}>
					{#if isLoading}
						<Icons.spinner class="mr-2 h-4 w-4 animate-spin" />
					{:else}
						<Icons.google class="mr-2 h-4 w-4" />
					{/if}
					{m.sign_in_with_google()}
				</Button>
			</div>

			<p class="text-muted-foreground px-8 text-center text-sm">
				{m.by_clicking_continue()}
				<a href="/terms" class="hover:text-primary underline underline-offset-4">
					{m.terms_of_service()}
				</a>
				and
				<a href="/privacy" class="hover:text-primary underline underline-offset-4">
					{m.privacy_policy()}
				</a>
				.
			</p>
		</div>
	</div>
</div>