import LeagueDropdown from "@/components/League/LeagueDropdown";
import {View, StyleSheet} from "react-native";
import {Button} from "react-native-paper";

interface SeasonsSearchBarProps {
    onLeagueSelection: (value: string) => void
}
export default function SeasonSearchBar(props: SeasonsSearchBarProps) {
    const oncChange = (value: string): void => {
        props.onLeagueSelection(value);
    }
    return (
        <View>
            <LeagueDropdown onChange={oncChange}></LeagueDropdown>
        </View>
    )
}

const styles = StyleSheet.create({

})
