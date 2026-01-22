import { writable, derived } from "svelte/store";
import type { PeriodResponse, CreatePeriodRequest, CreatePeriodsRequest, UpdatePeriodRequest } from "../dtos/period-dto";
import {
    getPeriods,
    getPeriod,
    createPeriod,
    createPeriods,
    updatePeriod,
    deletePeriod,
} from "../services/period-service";

function createPeriodStore() {
    // ===== state =====
    const periods = writable<PeriodResponse[]>([]);
    const selectedPeriod = writable<PeriodResponse | null>(null);
    const searchTerm = writable("");
    const currentPage = writable(1);
    const itemsPerPage = writable(5);

    // ===== derived =====
    const filteredPeriods = derived(
        [periods, searchTerm],
        ([$periods, $search]) =>
            $periods.filter((p) =>
                p.title.toLowerCase().includes($search.toLowerCase())
            )
    );

    const totalPages = derived(
        [filteredPeriods, itemsPerPage],
        ([$filtered, $limit]) =>
            Math.max(1, Math.ceil($filtered.length / $limit))
    );

    const paginatedPeriods = derived(
        [filteredPeriods, currentPage, itemsPerPage],
        ([$filtered, $page, $limit]) => {
            const start = ($page - 1) * $limit;
            return $filtered.slice(start, start + $limit);
        }
    );

    // ===== actions =====
    function setData(data: PeriodResponse[]) {
        periods.set(data);
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

    async function fetchPeriods(classroomId: string, roundId: string) {
        try {
            const data = await getPeriods(classroomId, roundId);
            setData(data);
        } catch (err) {
            console.error("Failed to fetch periods", err);
            setData([]);
        }
    }

    async function fetchPeriod(classroomId: string, periodId: string) {
        try {
            const data = await getPeriod(classroomId, periodId);
            selectedPeriod.set(data);
        } catch (err) {
            console.error("Failed to fetch period", err);
            selectedPeriod.set(null);
        }
    }

    async function addPeriod(data: CreatePeriodRequest, classroomId: string, roundId: string) {
        try {
            const newPeriod = await createPeriod(data, classroomId, roundId);
            if (!newPeriod || !newPeriod.id) {
                console.error("Dữ liệu trả về từ Service bị sai cấu trúc:", newPeriod);
                throw new Error("Invalid response structure");
            }
            periods.update((list) => [...list, newPeriod]);
            return newPeriod;
        } catch (err) {
            console.error("Failed to create period", err);
            throw err;
        }
    }

    async function addPeriods(data: CreatePeriodsRequest) {
        try {
            const result = await createPeriods(data);
            await fetchPeriods(data.classroomId, data.roundId); // Refresh list
            return result;
        } catch (err) {
            console.error("Failed to create periods", err);
            throw err;
        }
    }

    async function updatePeriodData(data: UpdatePeriodRequest) {
        try {
            const updatedPeriod = await updatePeriod(data);
            periods.update((list) =>
                list.map((p) => (p.id === data.periodId ? updatedPeriod : p))
            );
        } catch (err) {
            console.error("Failed to update period", err);
            throw err;
        }
    }

    async function removePeriod(classroomId: string, roundId: string, periodId: string) {
        try {
            await deletePeriod(classroomId, roundId, periodId);
            periods.update((list) => list.filter((p) => p.id !== periodId));
        } catch (err) {
            console.error("Failed to delete period", err);
            throw err;
        }
    }

    return {
        // state
        periods,
        selectedPeriod,
        searchTerm,
        currentPage,
        itemsPerPage,

        // derived
        filteredPeriods,
        paginatedPeriods,
        totalPages,

        // actions
        setData,
        setSearch,
        changePage,
        changePageSize,
        fetchPeriods,
        fetchPeriod,
        addPeriod,
        addPeriods,
        updatePeriodData,
        removePeriod,
    };
}

export const periodStore = createPeriodStore();
