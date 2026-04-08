import {Divider} from "react-native-paper";
import {View, StyleSheet} from "react-native";
import { AppButton } from "@/components/ui/AppButton";
import { AppCard } from "@/components/ui/AppCard";
import { AppText } from "@/components/ui/AppText";

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
    return (
        <View style={[styles.container]}>
            <AppCard style={[styles.card]}>
                <AppCard.Title title={"League Members"}/>
                <AppCard.Content>
                    <View>
                        {
                            props.members.map((member: LeagueMember) => {
                                return (
                                    <View key={member.id}>
                                        <Divider style={{marginTop: 4}}/>
                                        <View style={styles.row} key={member.id}>
                                            <View style={{flex: 1}}><AppText>{member.teamName}</AppText></View>
                                            <View>
                                                <AppButton mode={'contained'} style={[styles.button]} onPress={() => props.onRemove(member.id)}>Remove</AppButton>
                                            </View>
                                        </View>
                                    </View>
                                )
                            })
                        }
                    </View>
                </AppCard.Content>
            </AppCard>
        </View>
    )
}

const styles = StyleSheet.create({
    button: {
        borderRadius: 0,
        alignSelf: "flex-end"
    },
    row: {
        marginTop: 6,
        flex: 1,
        flexDirection: "row",
        justifyContent:"space-between",
        alignItems: "center"
    },
    container: {
        height: '50%',
        marginTop: 16
    },
    card: {
    }
});
