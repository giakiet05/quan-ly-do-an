// src/stores/class.store.ts
import { writable, derived } from "svelte/store";
import type { ClassItem } from "../types/class";



function createClassStore() {
    // ===== state =====
    const classes = writable<ClassItem[]>([]);
    const searchTerm = writable("");
    const currentPage = writable(1);
    const itemsPerPage = writable(5);

    // ===== derived =====
    const filteredClasses = derived(
        [classes, searchTerm],
        ([$classes, $search]) =>
            $classes.filter(
                (c) =>
                    c.name.toLowerCase().includes($search.toLowerCase()) ||
                    c.semester.toLowerCase().includes($search.toLowerCase())
            )
    );

    const totalPages = derived(
        [filteredClasses, itemsPerPage],
        ([$filtered, $limit]) =>
            Math.max(1, Math.ceil($filtered.length / $limit))
    );

    const paginatedClasses = derived(
        [filteredClasses, currentPage, itemsPerPage],
        ([$filtered, $page, $limit]) => {
            const start = ($page - 1) * $limit;
            return $filtered.slice(start, start + $limit);
        }
    );

    // ===== actions =====
    function setData(data: ClassItem[]) {
        classes.set(data);
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

    function addClass(item: ClassItem) {
        classes.update((list) => [...list, item]);
    }

    function updateClass(item: ClassItem) {
        classes.update((list) =>
            list.map((c) => (c.name === item.name ? item : c))
        );
    }

    function removeClass(name: string) {
        classes.update((list) => list.filter((c) => c.name !== name));
    }

    return {
        // state
        classes,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredClasses,
        paginatedClasses,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        addClass,
        updateClass,
        removeClass,
    };
}

export const classStore = createClassStore();
