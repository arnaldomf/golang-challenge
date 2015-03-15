package drum

import "os"
import "io"
import "encoding/binary"
import "strings"
import "fmt"
import "errors"
// DecodeFile decodes the drum machine file found at the provided path
// and returns a pointer to a parsed pattern which is the entry point to the
// rest of the data.
// TODO: implement
func DecodeFile(path string) (*Pattern, error) {
	p 			:= &Pattern{}
	pfile, err  := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return p, err
	}

	defer pfile.Close()

	err = ReadHeader(pfile, p)
	if err != nil {
		return p, err
	}
	// tracks comecam no byte 50 = Size - 36
	
	err, tracks := ReadData(pfile, p)
	if err != nil {
		return p, err
	}
	p.Tracks = tracks
	return p, err
}

func ReadHeader(file *os.File, p *Pattern) error {
	var header Header
	//Read file and put on header
	err := binary.Read(file, binary.LittleEndian, &header)
	if err == nil {
		if ! strings.Contains("SPLICE", string(header.Magic[:])) {
			return errors.New("Not SPLICE")
		}
		p.Header = header
	}
	return err
}

func ReadData(file *os.File, p *Pattern) (error, []Track) {
	data := make([]byte, p.Size - 36)
	_, err := io.ReadFull(file, data)
	i := 0
	var tracks []Track
	for i < len(data) {
		var track Track
		fmt.Println(data[i:i+4])
		x, n := binary.Uvarint(data[i:i+4])
		fmt.Println(n)
		if n < 0 {
			return err, tracks
		}
		track.Id = uint32(x)
		i += 4 // li mais 4
		track.NameSize = data[i]
		i += 1
		track.Name = data[i:int(track.NameSize) + i]
		i += int(track.NameSize)
		
		for j := i;j < i + 16; j++ {
			if (j - i) % 4 == 0 {
				track.Steps = append(track.Steps,byte(0x7c))
			}
			if data[j] == 0 {
				track.Steps = append(track.Steps,byte(0x2d))
			} else {
				track.Steps = append(track.Steps,byte(0x78))
			}
		}
		i += 16
		tracks = append(tracks, track)
	}
	return err, tracks
}

// Pattern is the high level representation of the
// drum pattern contained in a .splice file.
// TODO: implement
type Header struct {
	Magic   [6]byte
	_       [7]byte
	Size    byte
	Version [32]byte
	Tempo	float32
}
type Track struct {
	Id 	  		uint32
	NameSize  	byte
	Name  		[]byte
	Steps 		[]byte
}
type Pattern struct{
	Header
	Tracks  []Track
}

func (p *Pattern) String() string {
	v := string(p.Version[:])
	t := p.Tempo
	var tracks string
	track_tmpl := "(%d) %s\t%s|\n"
	for _,t := range p.Tracks {
		x := string(t.Steps[:])
		tracks += fmt.Sprintf(track_tmpl, int8(t.Id), t.Name, x)
	}
	return fmt.Sprintf("Saved with HW Version: %s\nTempo: %v\n%s",
		strings.Trim(v,"\x00"), t, tracks)

}

// func main() {

// 	_, err := DecodeFile("fixtures/pattern_1.splice")
// 	fmt.Println(err)
// }