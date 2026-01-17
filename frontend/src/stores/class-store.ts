// src/stores/class.store.ts
import { writable, derived } from "svelte/store";
import type { ClassItem, CreateClassRequest } from "../types/class";
import { createClassroom, getMyClassrooms, updateClassroom } from "../services/classroom-service";
import { mapClassroomToClassItem, mapUIRequestToDTO } from "../mappers/classroom-mapper";
import type { ClassroomResponse } from "../dtos/classroom-dto";



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

    async function addClass(uiData: CreateClassRequest) {
        try {
            const dtoPayload = mapUIRequestToDTO(uiData);
            const newClassRaw = await createClassroom(dtoPayload);
            console.log("New class created:", newClassRaw);
            const newClassMapped = mapClassroomToClassItem(newClassRaw);
            classes.update((list) => [newClassMapped, ...list]);

            return newClassMapped;
        } catch (err) {
            console.error("Failed to create classroom", err);
            throw err;
        }
    }

    async function updateClass(id: string, uiData: CreateClassRequest) {
        try {
            const dtoPayload = mapUIRequestToDTO(uiData);
            const response = await updateClassroom(id, dtoPayload);

            console.log("API update response:", response);

            await fetchMyClasses();
        } catch (err) {
            console.error("Failed to update classroom", err);
            throw err;
        }
    }

    function removeClass(name: string) {
        classes.update((list) => list.filter((c) => c.name !== name));
    }

    async function fetchMyClasses() {
        try {
            const data = await getMyClassrooms();
            console.log("Dữ liệu nhận được:", data);
            const rawClasses = (data as any).classrooms || [];
            const mapped = rawClasses.map(mapClassroomToClassItem);
            setData(mapped);

        } catch (err) {
            console.error("Failed to fetch classrooms", err);
            setData([]); // Lỗi thì cho danh sách trống để không vỡ giao diện
        }
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
        fetchMyClasses,
    };
}

export const classStore = createClassStore();
