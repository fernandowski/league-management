import {apiRequest} from "@/api/api";
import {create} from 'zustand'
import {fetchJWT} from "@/util/jwt-manager";

interface OrganizationStore {
    organizations: any[];
    organization: string | null,
    loading: boolean;
    error: string | null;
    fetchOrganizations: () => void;
    setOrganization: (organizationId: string) => void;
}

export const useOrganizationStore = create<OrganizationStore>((set) => ({
    organizations: [],
    organization: null,
    loading: false,
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

            set({organizations: response, loading: false});

        } catch (error: any) {
            set({error: error.message, loading: false});
        }
    },
    setOrganization:  (organizationId: string) => {
        set({organization: organizationId})
    }
}));
