import {useCallback, useEffect, useState} from "react";
import LeagueMembers, {LeagueMember} from "@/components/League/LeagueMembers";
import {LeagueMembershipResponse, useData} from "@/hooks/useData";
import {StyleSheet, View, useWindowDimensions} from "react-native";
import LeagueMemberSearch from "@/components/League/LeagueMemberSearch";
import {useOrganizationStore} from "@/stores/organizationStore";
import {apiRequest} from "@/api/api";
import { AppCard } from "@/components/ui/AppCard";
import { AppButton } from "@/components/ui/AppButton";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";
import ControlledTextInput from "@/components/FormControls/ControlledTextInput";
import StyledModal from "@/components/StyledModal";
import Joi from "joi";
import {useForm} from "react-hook-form";
import {joiResolver} from "@hookform/resolvers/joi";

export interface MembershipManagementProps {
    leagueId: string | null
    onFetch?: () => void
}

interface CreateTeamData {
    name: string
}

const schema = Joi.object({
    name: Joi.string()
        .trim()
        .required()
        .messages({
            "string.empty": "Team name is required",
        })
});

export default function MembershipManagement(props: MembershipManagementProps) {
    const [leagueMembers, setLeagueMembers] = useState<LeagueMember[]>([]);
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [draftTeamName, setDraftTeamName] = useState("");
    const [createError, setCreateError] = useState<string | null>(null);
    const [creatingTeam, setCreatingTeam] = useState(false);
    const {fetchData} = useData<LeagueMembershipResponse[]>();
    const {organization} = useOrganizationStore();
    const theme = useAppTheme();
    const dimensions = useWindowDimensions();
    const isLargeScreen = dimensions.width >= 960;
    const {control, handleSubmit, reset, formState: {errors}} = useForm<CreateTeamData>({
        resolver: joiResolver(schema),
        defaultValues: {
            name: "",
        }
    });

    const fetchMembership = useCallback(async (notifyParent = false) => {
        if (!props.leagueId) {
            setLeagueMembers([]);
            return;
        }
        const response = await fetchData(`/v1/leagues/${props.leagueId}/members`);
        if (response) {
            const nextMembers: LeagueMember[] = response.map((member: LeagueMembershipResponse) => ({
                leagueId: member.league_id,
                teamName: member.team_name,
                teamId: member.team_id,
                id: member.membership_id,
            }));
            setLeagueMembers(nextMembers);
            if (notifyParent) {
                props.onFetch && props.onFetch()
            }
        }
    }, [fetchData, props.leagueId, props.onFetch])

    const addTeamToLeague = useCallback(async (teamId: string) => {
        await apiRequest(
            `/v1/leagues/${props.leagueId}/invites`,
            {
                method: 'POST',
                body: {team_id: teamId}
            },
        );
    }, [props.leagueId]);

    const handleCreateTeam = async (formData: CreateTeamData) => {
        if (!organization || !props.leagueId) {
            return;
        }

        setCreatingTeam(true);
        setCreateError(null);

        try {
            const trimmedName = formData.name.trim();

            await apiRequest('/v1/teams', {
                method: 'POST',
                body: {
                    name: trimmedName,
                    organization_id: organization
                }
            });

            const matchingTeams = await apiRequest(
                `/v1/teams/?organization_id=${organization}&term=${encodeURIComponent(trimmedName)}`,
                {method: 'GET'}
            ) as Array<{id: string; name: string}>;

            const nextTeam = matchingTeams.find((team) =>
                team.name.toLowerCase() === trimmedName.toLowerCase() &&
                !leagueMembers.some((member) => member.teamId === team.id)
            ) ?? matchingTeams.find((team) => !leagueMembers.some((member) => member.teamId === team.id));

            if (!nextTeam) {
                throw new Error("Team was created but could not be added to the league automatically.");
            }

            await addTeamToLeague(nextTeam.id);
            reset({name: ""});
            setShowCreateModal(false);
            await fetchMembership(true);
        } catch (error: any) {
            setCreateError(error?.message || "Unable to create team.");
        } finally {
            setCreatingTeam(false);
        }
    };

    const handleOnRemove = async (membershipId: string): Promise<void> => {
        await apiRequest(`/v1/leagues/${props.leagueId}/members/${membershipId}`, {method: 'DELETE'})
        fetchMembership(true);
    }

    useEffect(() => {
        fetchMembership();
    }, [fetchMembership]);

    useEffect(() => {
        if (showCreateModal) {
            reset({name: draftTeamName});
        }
    }, [draftTeamName, reset, showCreateModal]);

    const handleOpenCreateModal = (teamName: string) => {
        setCreateError(null);
        setDraftTeamName(teamName.trim());
        setShowCreateModal(true);
    };

    const handleCloseCreateModal = () => {
        setCreateError(null);
        setShowCreateModal(false);
    };

    return (
        <View style={styles.container}>
            <View style={styles.header}>
                <View style={styles.headerCopy}>
                    <AppText variant="titleMedium">Teams</AppText>
                    <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                        Create a new team or find an existing organization team, then keep the current league teams up to date.
                    </AppText>
                </View>
                <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                    {leagueMembers.length} team{leagueMembers.length === 1 ? "" : "s"}
                </AppText>
            </View>

            <View style={[styles.grid, isLargeScreen && styles.gridLarge]}>
                <View style={styles.actionsColumn}>
                <AppCard style={[styles.searchPanel, {borderColor: theme.colors.outline}]}>
                    <AppCard.Content style={styles.panelContent}>
                        <View style={styles.sectionHeader}>
                            <AppText variant="titleMedium">Find team</AppText>
                            <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                Search teams already in the organization and add them to this league.
                            </AppText>
                        </View>
                        <LeagueMemberSearch
                            organizationId={organization ? organization : ''}
                            leagueId={props.leagueId ? props.leagueId : ''}
                            onTeamAdded={() => fetchMembership(true)}
                            members={leagueMembers}
                            onCreateTeam={handleOpenCreateModal}
                            creatingTeam={creatingTeam}
                            createError={createError}
                        />
                    </AppCard.Content>
                </AppCard>
                </View>

                <View style={styles.membersPanel}>
                    <LeagueMembers members={leagueMembers} onRemove={handleOnRemove}/>
                </View>
            </View>

            <StyledModal
                isOpen={showCreateModal}
                onDismiss={handleCloseCreateModal}
                contentContainerStyle={styles.modal}
            >
                <View style={styles.modalContent}>
                    <View style={styles.sectionHeader}>
                        <AppText variant="titleMedium">Create team</AppText>
                        <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                            Confirm the full team name before creating it and adding it to this league.
                        </AppText>
                    </View>
                    {createError ? (
                        <AppText variant="bodySmall" style={{color: theme.colors.error}}>
                            {createError}
                        </AppText>
                    ) : null}
                    <ControlledTextInput
                        label="Team name"
                        name="name"
                        control={control}
                        error={errors.name?.message}
                    />
                    <View style={styles.modalActions}>
                        <AppButton variant="secondary" onPress={handleCloseCreateModal}>Cancel</AppButton>
                        <AppButton
                            variant="submit"
                            onPress={handleSubmit(handleCreateTeam)}
                            loading={creatingTeam}
                            disabled={creatingTeam}
                        >
                            Create and add
                        </AppButton>
                    </View>
                </View>
            </StyledModal>
        </View>
    )
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        gap: 16,
    },
    header: {
        flexDirection: "row",
        flexWrap: "wrap",
        gap: 12,
        justifyContent: "space-between",
        alignItems: "center",
    },
    headerCopy: {
        flex: 1,
        minWidth: 240,
        gap: 4,
    },
    grid: {
        gap: 16,
    },
    gridLarge: {
        flexDirection: "row",
        alignItems: "flex-start",
    },
    actionsColumn: {
        flex: 1,
        gap: 16,
    },
    searchPanel: {
        borderRadius: 18,
    },
    panelContent: {
        gap: 14,
    },
    sectionHeader: {
        gap: 4,
    },
    membersPanel: {
        flex: 1.2,
    },
    modal: {
        width: "90%",
        maxWidth: 440,
        alignSelf: "center",
    },
    modalContent: {
        padding: 16,
        gap: 14,
    },
    modalActions: {
        flexDirection: "row",
        justifyContent: "flex-end",
        gap: 8,
    },
});
