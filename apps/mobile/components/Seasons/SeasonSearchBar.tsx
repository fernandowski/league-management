import LeagueDropdown from "@/components/League/LeagueDropdown";
import {View} from "react-native";

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
