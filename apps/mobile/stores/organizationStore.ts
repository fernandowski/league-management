import {apiRequest} from "@/api/api";
import {create} from 'zustand'
import {fetchJWT} from "@/util/jwt-manager";

interface OrganizationStore {
    organizations: Organization[];
    organization: string | null,
    loading: boolean;
    initialized: boolean;
    error: string | null;
    fetchOrganizations: () => void;
    setOrganization: (organizationId: string) => void;
}

export interface Organization {
    id: string;
    name: string;
}

export const useOrganizationStore = create<OrganizationStore>((set) => ({
    organizations: [],
    organization: null,
    loading: false,
    initialized: false,
    error: null,
    fetchOrganizations: async () => {
        set({loading: true});
        try {
            const jwt = await fetchJWT();
            const response = await apiRequest('/v1/organizations', {
                method: 'GET',
                headers: {
                    auth: jwt as string
                }
            });

            set((state) => {
                const nextOrganizations = response as Organization[];
                const hasSelectedOrganization = nextOrganizations.some(
                    ({id}) => id === state.organization
                );

                return {
                    organizations: nextOrganizations,
                    organization: hasSelectedOrganization
                        ? state.organization
                        : nextOrganizations[0]?.id ?? null,
                    loading: false,
                    initialized: true,
                    error: null,
                };
            });

        } catch (error: any) {
            set({error: error.message, loading: false, initialized: true});
        }
    },
    setOrganization:  (organizationId: string) => {
        set({organization: organizationId})
    }
}));
