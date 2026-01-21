import { writable, get } from "svelte/store";
import {
    getProjectRounds,
    createProjectRound,
    deleteProjectRound,
    updateProjectRound,
    getProjectRoundById
} from "../services/project-round-service";
import type { ProjectRound } from "../types/project-round";
import type { CreateProjectRoundRequest, UpdateProjectRoundRequest } from "../dtos/project-dto";

const roundsStore = writable<ProjectRound[]>([]);
const isLoadingStore = writable(false);
const errorStore = writable<string | null>(null);
const currentRoundStore = writable<ProjectRound | null>(null);

export const projectRoundStore = {
    // Để component có thể dùng $projectRoundStore
    subscribe: roundsStore.subscribe,
    get currentRounds() {
        return get(roundsStore);
    },

    async fetchRounds(classroomId: string) {
        isLoadingStore.set(true);
        try {
            const data = await getProjectRounds(classroomId);
            roundsStore.set(data);
        } catch (e) {
            errorStore.set("Không thể tải danh sách vòng dự án");
        } finally {
            isLoadingStore.set(false);
        }
    },

    async addRound(data: CreateProjectRoundRequest) {
        try {
            const newRound = await createProjectRound(data);
            roundsStore.update(n => [...n, newRound]);
            return newRound;
        } catch (e) {
            errorStore.set("Lỗi khi tạo vòng dự án");
            throw e;
        }
    },
    async editRound(data: UpdateProjectRoundRequest) {
        try {
            const updated = await updateProjectRound(data);
            roundsStore.update(list =>
                list.map(r => (r.id === data.roundId ? { ...r, ...updated } : r))
            );
        } catch (e) {
            errorStore.set("Lỗi khi cập nhật hạng mục");
            throw e;
        }
    },

    async removeRound(roundId: string, classroom_id: string) {

        try {
            await deleteProjectRound(roundId, classroom_id);
            roundsStore.update(n => n.filter(r => r.id !== roundId));
        } catch (e) {
            errorStore.set("Lỗi khi xóa vòng dự án");
            alert("Xóa thất bại!");
        }
    },
    async getRoundById(classroomId: string, roundId: string) {
        try {
            const data = await getProjectRoundById(classroomId, roundId);
            currentRoundStore.set(data); // Set the current round
            return data;
        } catch (e) {
            errorStore.set("Lỗi khi lấy chi tiết vòng dự án");
            throw e;
        }
    },
    get currentRound() {
        return get(currentRoundStore);
    }
};

// Export thêm các trạng thái phụ dưới dạng readable để an toàn
export const projectRoundLoading = { subscribe: isLoadingStore.subscribe };
export const projectRoundError = { subscribe: errorStore.set };
// Export the current round as a readable store
export const currentProjectRound = { subscribe: currentRoundStore.subscribe };