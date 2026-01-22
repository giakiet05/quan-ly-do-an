import { writable, derived } from "svelte/store";
import type { PostResponse, CreatePostRequest } from "../dtos/post-dto";
import {
    getPosts,
    getPost,
    createPost,
    updatePost,
    deletePost,
} from "../services/post-service";

function createPostStore() {
    // ===== state =====
    const posts = writable<PostResponse[]>([]);
    const selectedPost = writable<PostResponse | null>(null);
    const searchTerm = writable("");
    const currentPage = writable(1);
    const itemsPerPage = writable(5);

    // ===== derived =====
    const filteredPosts = derived(
        [posts, searchTerm],
        ([$posts, $search]) =>
            $posts.filter((p) =>
                p.title.toLowerCase().includes($search.toLowerCase())
            )
    );

    const totalPages = derived(
        [filteredPosts, itemsPerPage],
        ([$filtered, $limit]) =>
            Math.max(1, Math.ceil($filtered.length / $limit))
    );

    const paginatedPosts = derived(
        [filteredPosts, currentPage, itemsPerPage],
        ([$filtered, $page, $limit]) => {
            const start = ($page - 1) * $limit;
            return $filtered.slice(start, start + $limit);
        }
    );

    // ===== actions =====
    function setData(data: PostResponse[]) {
        posts.set(data);
    }

    function setSearch(value: string) {
        searchTerm.set(value);
        currentPage.set(1);
    }

    function changePage(page: number) {
        currentPage.set(page);
    }

    function changePageSize(size: number) {
        itemsPerPage.set(size);
        currentPage.set(1);
    }

    async function fetchPosts(classroomId: string) {
        try {
            const data = await getPosts(classroomId);
            console.log("Fetched posts:", data);
            setData(data);
        } catch (err) {
            console.error("Failed to fetch posts", err);
            setData([]);
        }
    }

    async function fetchPost(classroomId: string, postId: string) {
        try {
            const data = await getPost(classroomId, postId);
            selectedPost.set(data);
        } catch (err) {
            console.error("Failed to fetch post", err);
            selectedPost.set(null);
        }
    }

    async function addPost(data: CreatePostRequest, classroomId: string) {
        try {
            const newPost = await createPost(data, classroomId);
            if (!newPost || !newPost.id) {
                console.error("Dữ liệu trả về từ Service bị sai cấu trúc:", newPost);
                throw new Error("Invalid response structure");
            }
            posts.update((list) => [...list, newPost]);
            return newPost;
        } catch (err) {
            console.error("Failed to create post", err);
            throw err;
        }
    }

    async function updatePostData(postId: string, data: Partial<CreatePostRequest>) {
        try {
            const updatedPost = await updatePost(postId, data);
            posts.update((list) =>
                list.map((p) => (p.id === postId ? updatedPost : p))
            );
        } catch (err) {
            console.error("Failed to update post", err);
            throw err;
        }
    }

    async function removePost(classroomId: string, postId: string) {
        try {
            await deletePost(classroomId, postId);
            posts.update((list) => list.filter((p) => p.id !== postId));
        } catch (err) {
            console.error("Failed to delete post", err);
            throw err;
        }
    }

    return {
        // state
        posts,
        selectedPost,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredPosts,
        paginatedPosts,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        fetchPosts,
        fetchPost,
        addPost,
        updatePostData,
        removePost,
    };
}

export const postStore = createPostStore();
