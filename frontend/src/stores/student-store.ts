import { writable, derived } from "svelte/store";
import type { StudentInClass } from "../types/student";

function createStudentStore() {
    // ===== state =====
    const students = writable<StudentInClass[]>([]);
    const searchTerm = writable("");
    const currentPage = writable(1);
    const itemsPerPage = writable(5);

    // ===== derived =====
    const filteredStudents = derived(
        [students, searchTerm],
        ([$students, $search]) =>
            $students.filter(
                (s) =>
                    s.name.toLowerCase().includes($search.toLowerCase()) ||
                    s.studentCode.toLowerCase().includes($search.toLowerCase()) ||
                    s.email.toLowerCase().includes($search.toLowerCase())
            )
    );

    const totalPages = derived(
        [filteredStudents, itemsPerPage],
        ([$filtered, $limit]) =>
            Math.max(1, Math.ceil($filtered.length / $limit))
    );

    const paginatedStudents = derived(
        [filteredStudents, currentPage, itemsPerPage],
        ([$filtered, $page, $limit]) => {
            const start = ($page - 1) * $limit;
            return $filtered.slice(start, start + $limit);
        }
    );

    // ===== actions =====
    function setData(data: StudentInClass[]) {
        students.set(data);
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

    function addStudent(student: StudentInClass) {
        students.update((list) => [...list, student]);
    }

    function updateStudent(student: StudentInClass) {
        students.update((list) =>
            list.map((s) => (s.id === student.id ? student : s))
        );
    }

    function removeStudent(id: string) {
        students.update((list) => list.filter((s) => s.id !== id));
    }

    return {
        // state
        students,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredStudents,
        paginatedStudents,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        addStudent,
        updateStudent,
        removeStudent,
    };
}

export const studentStore = createStudentStore();
