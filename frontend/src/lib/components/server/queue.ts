const queue: { url: string, title: string, duration: string }[] = []

export const add_song = (song: { url: string, title: string, duration: string }) => {
    queue.push(song)
}

export const get_queue = () => {
    return queue
}

export const clear_queue = () => {
    queue.length = 0
}

export const remove_song = (song: { url: string, title: string, duration: string }) => {
    queue.splice(queue.indexOf(song), 1)
}