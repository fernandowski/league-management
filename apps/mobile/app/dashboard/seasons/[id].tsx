import {useLocalSearchParams, useRouter} from "expo-router";
import {useEffect, useRef, useState} from "react";
import SeasonManagement from "@/components/Seasons/SeasonManagement";
import ViewContent from "@/components/Layout/ViewContent";
import {useOrganizationStore} from "@/stores/organizationStore";


export default function SeasonDetails() {
    const {id: urlParamSeasonID} = useLocalSearchParams();
    const [seasonId, setSeasonID] = useState<string>(Array.isArray(urlParamSeasonID) ? urlParamSeasonID[0] : urlParamSeasonID || "");

    const {organization} = useOrganizationStore();
    const router = useRouter()
    const prevOrganizationRef = useRef<string | null>(organization);

    useEffect(() => {
        const newSeasonId = Array.isArray(urlParamSeasonID) ? urlParamSeasonID[0] : urlParamSeasonID || ""
        setSeasonID(newSeasonId)
    }, [urlParamSeasonID]);

    useEffect(() => {
        if (organization && prevOrganizationRef.current !== organization) {
            prevOrganizationRef.current = organization;
            router.push("/dashboard/seasons");
        }
    }, [organization, router]);

    return (
        <ViewContent>
            <SeasonManagement seasonId={seasonId}/>
        </ViewContent>)
}

