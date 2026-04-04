import {Button, Card, Divider, IconButton, Text} from "react-native-paper";
import {View, StyleSheet} from "react-native";

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
            <Card style={[styles.card]}>
                <Card.Title title={"League Members"}/>
                <Card.Content>
                    <View>
                        {
                            props.members.map((member: LeagueMember) => {
                                return (
                                    <View key={member.id}>
                                        <Divider style={{marginTop: 4}}/>
                                        <View style={styles.row} key={member.id}>
                                            <View style={{flex: 1}}><Text>{member.teamName}</Text></View>
                                            <View>
                                                <Button mode={'contained'} style={[styles.button]} onPress={() => props.onRemove(member.id)}>Remove</Button>
                                            </View>
                                        </View>
                                    </View>
                                )
                            })
                        }
                    </View>
                </Card.Content>
            </Card>
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
