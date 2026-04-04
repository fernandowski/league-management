import {useEffect, useState} from "react";
import LeagueMembers, {LeagueMember} from "@/components/League/LeagueMembers";
import {LeagueMembershipResponse, useData} from "@/hooks/useData";
import {View} from "react-native";
import LeagueMemberSearch from "@/components/League/LeagueMemberSearch";
import {useOrganizationStore} from "@/stores/organizationStore";
import {apiRequest} from "@/api/api";

export interface MembershipManagementProps {
    leagueId: string | null
    onFetch?: () => void
}

export default function MembershipManagement(props: MembershipManagementProps) {
    const [leagueMembers, setLeagueMembers] = useState<LeagueMember[]>([]);
    const {fetchData, fetching, error} = useData<LeagueMembershipResponse[]>();
    const {organization} = useOrganizationStore();

    const fetchMembership = async () => {
        const response = await fetchData(`/v1/leagues/${props.leagueId}/members`);
        if (response) {
            const leagueMembers: LeagueMember[] = response.map((member: LeagueMembershipResponse) => ({
                leagueId: member.league_id,
                teamName: member.team_name,
                teamId: member.team_id,
                id: member.membership_id,
            }));
            setLeagueMembers(leagueMembers);
            props.onFetch && props.onFetch()
        }
    }

    const handleOnRemove = async (membershipId: string): Promise<void> => {
        await apiRequest(`/v1/leagues/${props.leagueId}/members/${membershipId}`, {method: 'DELETE'})
        fetchMembership();
    }
    useEffect(() => {
        fetchMembership();
    }, [props.leagueId]);
    return (
        <View style={{flex: 1}}>
            <View style={{marginTop: 8}}>
                <LeagueMemberSearch
                    organizationId={organization ? organization : ''}
                    leagueId={props.leagueId ? props.leagueId : ''}
                    onTeamAdded={fetchMembership}
                    members={leagueMembers}
                />
            </View>
            <LeagueMembers members={leagueMembers} onRemove={handleOnRemove}/>
        </View>

    )
}
