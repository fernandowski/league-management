import {View, StyleSheet} from "react-native";
import Tabs from "@/components/Layout/Tabs";
import SeasonStanding from "@/components/Seasons/SeasonStanding";
import SeasonMatchUpManagement from "@/components/Seasons/SeasonMatchUpManagement";


export interface SeasonManagementProps {
    seasonId: string
}

export default function SeasonManagement(props: SeasonManagementProps) {
    return (
        <View>
            <View>
                <Tabs tabs={[
                    {key: 'standing', title: 'Standings', view: <SeasonStanding seasonId={props.seasonId}/>},
                    {key: 'match-ups', title: 'Match Ups', view: <SeasonMatchUpManagement seasonId={props.seasonId}/>}
                ]}
                />
            </View>
        </View>
    )
}

const styles = StyleSheet.create({})
