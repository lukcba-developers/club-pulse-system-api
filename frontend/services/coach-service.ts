// Removed unused api import

export interface TravelEvent {
    id: string;
    title: string;
    destination: string;
    departure_date: string;
    meeting_time: string;
    type: 'TRAVEL' | 'MATCH' | 'TOURNAMENT' | 'TRAINING';
}

export const coachService = {
    getTravelEvents: async (): Promise<TravelEvent[]> => {
        // Mock implementation or future API endpoint
        return [];
    }
};
