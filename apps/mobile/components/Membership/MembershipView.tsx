import {View} from "react-native";
import MembershipManagement from "@/components/League/MembershipManagement";

export const MembershipView = (props: { leagueId: string, onMemberRefresh: () => void }) => {
    return (
        <View style={{flex: 1}}>
            <MembershipManagement leagueId={props.leagueId} onFetch={props.onMemberRefresh}/>
        </View>
    )
}

export default MembershipView
