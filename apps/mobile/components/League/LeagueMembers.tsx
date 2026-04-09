import {View, StyleSheet} from "react-native";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";
import { useAppTheme } from "@/theme/theme";

export interface LeagueMember {
    teamName: string
    id: string
    teamId: string
    leagueId: string
}

export interface LeagueMembersProps {
    members: LeagueMember[]
    onRemove: (id: string) => void
}

export default function LeagueMembers(props: LeagueMembersProps) {
    const theme = useAppTheme();

    return (
        <View style={[styles.container]}>
            <AppCard style={[styles.card]}>
                <AppCard.Content style={styles.content}>
                    <View style={styles.header}>
                        <View style={styles.headerCopy}>
                            <AppText variant="titleMedium">League members</AppText>
                            <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                Review the current lineup and remove teams when needed.
                            </AppText>
                        </View>
                        <AppText variant="labelLarge" style={{color: theme.colors.onSurfaceVariant}}>
                            {props.members.length} team{props.members.length === 1 ? "" : "s"}
                        </AppText>
                    </View>

                    {props.members.length === 0 ? (
                        <View style={[styles.emptyState, {backgroundColor: theme.colors.surfaceVariant, borderColor: theme.colors.outlineVariant}]}>
                            <AppText variant="bodyMedium" style={{color: theme.colors.onSurfaceVariant}}>
                                No teams have been added to this league yet.
                            </AppText>
                        </View>
                    ) : (
                        <View style={styles.list}>
                            {
                                props.members.map((member: LeagueMember) => {
                                    return (
                                        <View
                                            style={[styles.row, {borderColor: theme.colors.outlineVariant, backgroundColor: theme.colors.surfaceVariant}]}
                                            key={member.id}
                                        >
                                            <View style={styles.rowContent}>
                                                <AppText variant="titleMedium">{member.teamName}</AppText>
                                                <AppText variant="bodySmall" style={{color: theme.colors.onSurfaceVariant}}>
                                                    Team ID: {member.teamId}
                                                </AppText>
                                            </View>
                                            <View>
                                                <AppButton mode={'contained-tonal'} style={[styles.button]} onPress={() => props.onRemove(member.id)}>Remove</AppButton>
                                            </View>
                                        </View>
                                    )
                                })
                            }
                        </View>
                    )}
                </AppCard.Content>
            </AppCard>
        </View>
    )
}

const styles = StyleSheet.create({
    button: {
        alignSelf: "flex-end"
    },
    row: {
        borderWidth: 1,
        borderRadius: 18,
        padding: 14,
        flexDirection: "row",
        justifyContent:"space-between",
        alignItems: "center",
        gap: 12,
    },
    container: {
        marginTop: 4,
        flex: 1,
    },
    card: {
        borderRadius: 24,
    },
    content: {
        gap: 14,
    },
    header: {
        flexDirection: "row",
        justifyContent: "space-between",
        gap: 12,
        flexWrap: "wrap",
        alignItems: "center",
    },
    headerCopy: {
        flex: 1,
        gap: 4,
        minWidth: 220,
    },
    list: {
        gap: 10,
    },
    rowContent: {
        flex: 1,
        gap: 2,
    },
    emptyState: {
        borderWidth: 1,
        borderRadius: 18,
        padding: 16,
    },
});
